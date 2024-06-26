import React, { Component } from "react";
import Form from "react-validation/build/form";

export default function Login({ onEnterCode, error }) {
  return (
    <div className="background center-container">
      <div className="inner-content px-5">
        <div className="row">
          <div className="col-12">
            <div className="text-center mb-4">
              <h1>{process.env.REACT_APP_SCHOOL}</h1>
              <h1 className="h3 mb-3 font-weight-normal">Сканируйте QR код</h1>
            </div>

            <div class="mt-2">
              <input placeholder="Введите код подключения" onChange={onEnterCode} style={{ width: "100%"}}></input>
            </div>

            {error ? (
              <div className="alert alert-danger text-center" role="alert">
                {error}
              </div>
            ) : (<div></div>)}

            <p className="mt-5 mb-3 text-muted text-center">ALEGOR © 2023</p>
          </div>
        </div>
      </div>
    </div>
  );
}