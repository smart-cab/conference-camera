import React, { Component } from "react";
import Form from "react-validation/build/form";
import styles from "./login.css";
import QrReader from 'react-qr-scanner'

import { withRouter } from '../../common/with-router';

class Login extends Component {
  constructor(props) {
    super(props)
    this.state = {
      delay: 100,
      result: 'No result',
    }

    this.handleScan = this.handleScan.bind(this)
  }
  handleScan(data){
    if (data) {
      this.setState({
        result: data.text,
      })
    }
  }
  handleError(err) {
    console.error(err)
  }


  render() {
    const previewStyle = {
      width: 390,
      "border-radius": 10,
    }

    return (
      <div className="background">
        <Form
          onSubmit={this.handleLogin}
          ref={c => {
            this.form = c;
          }}
          className="formSignin"
        >
          <div className="text-center mb-4">
            <h1>1357</h1>
            <h1 className="h3 mb-3 font-weight-normal">Сканируйте QR код</h1>
          </div>

          <QrReader
            delay={this.state.delay}
            onError={this.handleError}
            style={previewStyle}
            onScan={this.handleScan}
          />

          <div className="form-group">
            <div className="alert alert-danger" role="alert">
              {this.state.result}
            </div>
          </div>

          <p className="mt-5 mb-3 text-muted text-center">ALEGOR © 2023</p>
        </Form>
      </div>
    );
  }
}

export default withRouter(Login);