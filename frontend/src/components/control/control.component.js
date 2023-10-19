import React, { Component } from "react";

import UserService from "../../services/user.service";
import styles from "./control.css";
import authService from "../../services/auth.service";
import Swal from 'sweetalert2'

export default class Control extends Component {
  constructor(props) {
    super(props);

    this.state = {
      user: authService.getCurrentUser(),
    };
  }

  componentDidMount() {
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
      <div>
        CAMERA FRAME
        <button>CENTER</button>
      </div>
    );
  }
}
