import React, { Component } from "react";
import {QRCodeSVG} from 'qrcode.react'

export default function Qr({ token }) {
    const tokenURL = `http://${window.location.hostname}:3000/?token=${token}`;

    return (
        <div className="container h-100">
          <div className="row h-100 justify-content-center align-items-center mt-5">
            <div className="col-12 text-center">
              <h1>{process.env.REACT_APP_SCHOOL} — Камера конференции</h1>
              <div className="mt-5">
                <QRCodeSVG value={token} size={256} />
              </div>
              <div className="mt-5">
                <a href={tokenURL} style={{ fontSize: "32px" }}>{token}</a>
              </div>
            </div>
          </div>
        </div>
    );
}