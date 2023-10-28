import React, { Component } from "react";
import {QRCodeSVG} from 'qrcode.react'

export default function Qr({ token }) {
    return (
        <div className="container h-100">
          <div className="row h-100 justify-content-center align-items-center mt-5">
            <div className="col-12 text-center">
              <h1>1234 — Камера конференции</h1>
              <div className="mt-5">
                <QRCodeSVG value={token} size={256} />
              </div>
              <div className="mt-5">
                <div>DEBUG INFO</div>
                <a href={token}>{token}</a>
              </div>
            </div>
          </div>
        </div>
    );
}