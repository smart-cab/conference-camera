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

    this.handleScan = this.handleScan.bind(this)
    this.socket = new WebSocket('wss://192.168.1.13:8888/ws');
  }


  componentDidMount() {
    this.socket.onopen = () => {
      this.socket.send("user:init");
      this.setState({ connecting: true })
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
      }
    };

    this.socket.onclose = (event) => {
      
    };

    this.socket.onerror = (error) => {
      
    };
  }

  handleScan(data){
    if (data) {
      // отправляем на бек токен с qr кода
      this.socket.send("user:connect:" + data.text)
    }
  }
  handleError(err) {
    console.error(err)
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