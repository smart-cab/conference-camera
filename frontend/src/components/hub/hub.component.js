import React, { Component } from "react";

import styles from "./hub.css";
import Connection from '../connection.component';
import Qr from "./qr.component";
import Error from "./error.component";
import Waiting from "./waiting.component";

export default class Hub extends Component {
  constructor(props) {
    super(props);

    this.state = {
      token: '',
      error: null,
      connecting: false,
      user: null,
      autoqr: false,
    };

    this.socket = new WebSocket(`ws://${window.location.hostname}:8888/api/v1/ws`);
  }

  componentDidMount() {
    this.socket.onopen = () => {
      // Отсылаем на бек о том что запущен hub
      this.socket.send("hub:init");
      this.setState({ connecting: true });
    };

    this.socket.onmessage = (event) => {
      const response = event.data;
      console.log('Received response from socket:', response);
      if (response.startsWith("token:")) {
        // Получаем токен от бека и генерируем новый QR код
        const token = `http://${window.location.hostname}:3000/?token=${response.replace("token:", "")}`;
        this.setState({ token });
      } 
      else if (response.startsWith("error:")) 
      {
        // Если получена ошибка
        const error = response.replace("error:", "");
        this.socket.close();
        if (error == 'already') {
          this.setState({ error: "Хаб уже запущен на другом устройстве." });
        } else {
          this.setState({ error: "Неизвестная ошибка." });
        }
      } 
      else if (response.startsWith("connected:")) 
      {
        // Если пользователь подключился
        this.setState({ user: response.replace("connected:", "") })
      } 
      else if (response == "disconnected") 
      {
        // Если пользователь отключился
        this.setState({ user: null })
      } 
      else if (response.startsWith("autoqr:")) 
      {
        // Проверяем включена ли на беке автогенерация QR кода
        if (response.replace("autoqr:", "") == "1") {
          this.setState({ autoqr: true })
          this.generateQRCodeValue(true);
        }
      }
    };

    this.socket.onclose = (event) => {
      this.setState({ connecting: false });
    };

    this.socket.onerror = (error) => {
      this.setState({ error: error.message, connecting: false });
    };

    setInterval(() => {
      this.generateQRCodeValue(false);
    }, 5 * 1000);
  }

  generateQRCodeValue(force) {
    if ((!this.state.user && !this.state.error && this.state.autoqr) || force) {
      this.socket.send("hub:token");
    } 
  }

  render() {
    if (!this.state.connecting) {
      // Нет подключения к сокету
      return (
        <Connection />
      );
    }

    if (this.state.error) {
      // Возникла ошибка
      return (
        <Error error={this.state.error} />
      );
    }

    if (this.state.user) {
      // Если пользователь авторизован, то выводим заставку
      return (
        <Waiting user={this.state.user}/>
      )
    }

    if (!this.state.token) {
      // Если отключена автоматическая генерация qr кода 
      return (
        <Error error="Ожидание запроса на подключение..." />
      );
    }

    return (
      // Отображаем QR код для подключения
      <Qr token={this.state.token} />
    );
  }
}
