import React, { Component } from "react";
import Form from "react-validation/build/form";
import QrReader from 'react-qr-scanner'

export default function Control() {
  const videoURL = `http://${window.location.hostname}:8888/api/v1/video`;
  return (
    <div className="background">
      <h1>CONTROLLER</h1>
      
      <img src={videoURL} />
    </div>
  );
}