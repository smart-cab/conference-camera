import React, { Component } from "react";

import UserService from "../../services/user.service";
import styles from "./hub.css";
import hubService from "../../services/hub.service";
import Swal from 'sweetalert2'
import {QRCodeSVG} from 'qrcode.react'
import md5 from 'md5'; // Импортируем библиотеку для хеширования

export default class Hub extends Component {
  constructor(props) {
    super(props);

    this.state = {
      timestamp: new Date().getTime(),
      secret: 'opcfsidr929320szas',
      token: ''
    };
  }

  generateQRCodeValue() {
    const { secret } = this.state;
    const currentTimeInSeconds = Math.floor(this.state.timestamp / 1000);
    const combinedString = `${currentTimeInSeconds}${secret}`;
    const hashedString = md5(combinedString);
    return hashedString;
  }

  componentDidMount() {
    setInterval(() => {
      this.setState({ timestamp: new Date().getTime() });
    }, 30 * 1000);

    const socket = new WebSocket('ws://127.0.0.1:8768/ws');


    socket.onopen = function () {
      console.log('Соединение установлено');
    };

    socket.onmessage = (event) => {
      const parsedData = JSON.parse(event.data);

    };

    socket.onclose = function (event) {
      console.log('Соединение закрыто');
    };

    socket.onerror = function (error) {
      console.log(`Ошибка: ${error.message}`);
    };
  }

  render() {
    return (
      <div class="container h-100">
        <div class="row h-100 justify-content-center align-items-center mt-5">
          <div class="col-12 text-center">
            <h1>1357 — Камера конференции</h1>
            <div class="mt-5">
              <QRCodeSVG value={this.generateQRCodeValue()} size={256} />
            </div>
            <div class="mt-5">
              <div>DEBUG INFO</div>
              {this.state.timestamp} — 192.168.1.1 — {this.generateQRCodeValue()}
            </div>
          </div>
        </div>
      </div>
    );
  }
}
