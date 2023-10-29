import React from "react";
import styles from "./control.css";

export default function Control() {
  const videoURL = `http://${window.location.hostname}:8888/api/v1/video`;
  return (
    <div className="container center-container">
      <div className="inner-content">
        <div className="row">
          <div className="col-md-7">
            <img src={videoURL} alt="Изображение" className="img-fluid video" />
          </div>
          <div className="col-md-5 text-center">
            <h3>1234 - Камера конференции</h3>
            <hr></hr>
            <div className="joystick-container">
              <div>
                <button className="joystick-button">↑</button>
              </div>
              <div>
                <button className="joystick-button">←</button>
                <button className="joystick-button">X</button>
                <button className="joystick-button">→</button>
              </div>
              <div>
                <button className="joystick-button">↓</button>
              </div>
            </div>
            <hr></hr>
            <span className="text-muted text-left">
              <ul>
                <li>DEBUG INFO:</li>
                <li>{videoURL}</li>
              </ul>
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
