import React, { Component } from "react";
import Form from "react-validation/build/form";

export default function Login({ onHandleScan, error }) {
  const previewStyle = {
    width: 390,
    borderRadius: 10,
  }

  return (
    <div className="background">
      <Form
        className="formSignin"
      >
        <div className="text-center mb-4">
          <h1>1234</h1>
          <h1 className="h3 mb-3 font-weight-normal">Сканируйте QR код</h1>
        </div>

        {error ? (
          <div className="form-group">
            <div className="alert alert-danger text-center" role="alert">
              {error}
           </div>
          </div>
        ) : (<div></div>)}

        <p className="mt-5 mb-3 text-muted text-center">ALEGOR © 2023</p>
      </Form>
    </div>
  );
}