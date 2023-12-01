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
      devices: {},
      selectedDevice: null,
      step: 300,
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
      } else if (response.startsWith("devices:")) {
        let devices = {}

        response.replace("devices:", "").split('|').forEach(device => {
          const [path, name] = device.split(':');
          if (path && name) {
            devices[name.trim()] = path.trim();
          }
        });

        this.setState({ devices: devices })

        // список девайсов камеры
      } else if (response.startsWith("selected-device:")) {
        let device = response.replace("selected-device:", "")
        console.log("selected device:", device)
        this.setState({ selectedDevice: device })
      }
    };

    this.socket.onclose = (event) => {

    };

    this.socket.onerror = (error) => {

    };
  }

  handleDeviceChange = (event) => {
    window.location.reload()
    this.socket.send("user:device:" + event.target.value);
  }

  moveCamera = (action) => {
    // alert("MOVE")
    console.log("move ptz:", action, ". Step:", this.state.step)
    this.socket.send("user:move:" + action + ":" + this.state.step)
  }

  zoomCamera = (event) => {
    console.log("zoom ptz:", event.target.value)
    this.socket.send("user:zoom:" + event.target.value)
  }

  stepSet = (event) => {
    this.setState({step : event.target.value})
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
        <Login onHandleScan={this.handleScan} error={this.state.error} />
      );
    }

    // Даем доступ к админке
    return (
      <Control devices={this.state.devices} selectedDevice={this.state.selectedDevice} deviceSelect={this.handleDeviceChange} moveCamera={this.moveCamera} zoomCamera={this.zoomCamera} stepSet={this.stepSet} />
    );
  }
}

export default withRouter(Main);