import React, { Component } from "react";
import styles from "./control.css";
import Login from "./login.component";
import Connection from '../connection.component';

import { withRouter } from '../../common/with-router';
import Control from "./control.component";

class Main extends Component {

  constructor(props) {
    super(props)
    this.state = {
      delay: 20000,
      auth: false,
      error: null,
      connecting: false,
    }

    console.log(window.location.hostname)
    this.socket = new WebSocket(`ws://${window.location.hostname}:8888/api/v1/ws`);

    let params = new URLSearchParams(window.location.search)
    this.token = params.get("token");
  }


  componentDidMount() {
    this.socket.onopen = () => {
      this.socket.send("user:init");
      this.setState({ connecting: true })
      this.socket.send("user:connect:" + this.token)
    };

    this.socket.onmessage = (event) => {
      const response = event.data;
      console.log('Received response from socket:', response);
      
      if (response == 'connected') {
        // токен успешно принят
        this.setState({ auth: true })
      } else if (response == 'disconnected') {
        // хаб отключился
        this.setState({ auth: false, connecting: false })
      } else if (response == 'already') {
        // уже есть подключенный пользователь
        this.setState({ error: 'Клиент уже подключен с другого устройства' })
      } else if (response == 'wrong') {
        this.setState({ error: 'Неверный токен авторизации' })
      }
    };

    this.socket.onclose = (event) => {
      
    };

    this.socket.onerror = (error) => {
      
    };
  }


  render() {
    if (!this.state.connecting) {
      // Нет подключения к сокету
      return (
        <Connection />
      );
    }

    if (!this.state.auth) {
      // Отображаем сканнер QR кода
      return (
        <Login onHandleScan={this.handleScan} error={this.state.error}/>
      );
    }

    // Даем доступ к админке
    return (
      <Control />
    );
  }
}

export default withRouter(Main);