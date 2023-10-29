import React from 'react';
import styles from "./control.css";

export default function Control({ devices, deviceSelect }) {
  const videoURL = `http://${window.location.hostname}:8888/api/v1/video`;

  return (
    <div className="container center-container">
      <div className="inner-content">
        <div className="row">
          <div className="col-md-7">
            <img src={videoURL} alt="Изображение" className="img-fluid video" />
          </div>
          <div className="col-md-5">
            <div className="text-center">
              <h3>1234 - Камера конференции</h3>
            </div>
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
            <h6>Настройки камеры</h6>
            <label className="pr-2">Камера:</label>
            <select onChange={deviceSelect}>
              {Object.entries(devices).map(([key, value]) => (
                <option key={key} value={value}>
                  {key}
                </option>
              ))}
            </select>
            <hr></hr>
            <span className="text-muted text-left">
              DEBUG INFO:
              <ul>
                <li>{videoURL}</li>
              </ul>
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
