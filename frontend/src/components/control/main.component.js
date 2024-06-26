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
      selectedScreen: null,
      step: 300,
      isPtz: false,
      image: '',
    }

    console.log(window.location.hostname)
    this.socket = new WebSocket(`ws://${window.location.hostname}:8888/api/v1/ws`);

    let params = new URLSearchParams(window.location.search)
    this.token = params.get("token");
  }


  componentDidMount() {
    this.enterCode = (event) => {
      this.socket.send("user:connect:" + event.target.value)
    }

    this.socket.onopen = () => {
      this.socket.send("user:init");
      this.setState({ connecting: true })
      if (this.token) {
        this.socket.send("user:connect:" + this.token)
      }
    };

    this.socket.onmessage = (event) => {
      const response = event.data;
      if (response.startsWith('data:image')) {
        this.setState({ image: response })
        return
      }

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
        const [deviceName, isPtz] = device.split(':')
        console.log("selected device:", deviceName, "ptz:", isPtz === "true")
        this.setState({ selectedDevice: deviceName, isPtz: isPtz === "true" })
      } else if (response.startsWith("selected-screen:")) {
        let device = response.replace("selected-screen:", "")
        const [deviceName, isPtz] = device.split(':')
        console.log("selected screen:", deviceName, "ptz:", isPtz === "true")
        this.setState({ selectedScreen: deviceName })
      }
    };

    this.socket.onclose = (event) => {

    };

    this.socket.onerror = (error) => {

    };
  }

  handleDeviceChange = (event) => {
    this.socket.send("user:switch:" + event.target.value);
  }

  handleScreenChange = (event) => {
    this.socket.send("user:dswitch:" + event.target.value);
  }

  moveCamera = (action) => {
    // alert("MOVE")
    console.log("move ptz:", action, ". Step:", this.state.step)
    this.socket.send("user:move:" + action + ":" + this.state.step)
  }

  changeScene = (action) => {
    // alert("MOVE")
    console.log("change scene:", action)
    this.socket.send("user:scene:" + action)
  }

  zoomCamera = (event) => {
    console.log("zoom ptz:", event.target.value)
    this.socket.send("user:zoom:" + event.target.value)
  }

  stepSet = (event) => {
    this.setState({step : event.target.value})
  }

  faceDetect = (event) => {
    this.socket.send("user:face:" + event.target.checked)
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
        <Login onEnterCode={this.enterCode} error={this.state.error} />
      );
    }

    // Даем доступ к админке
    return (
      <Control 
        devices={this.state.devices} 
        selectedDevice={this.state.selectedDevice} 
        deviceSelect={this.handleDeviceChange} 
        moveCamera={this.moveCamera} 
        changeScene={this.changeScene} 
        zoomCamera={this.zoomCamera} 
        stepSet={this.stepSet} 
        faceDetect={this.faceDetect}
        isPtz={this.state.isPtz}
        image={this.state.image}
        selectedScreen={this.state.selectedScreen} 
        screenSelect={this.handleScreenChange} 
      />
    );
  }
}

export default withRouter(Main);