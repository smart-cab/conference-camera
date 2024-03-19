import React from 'react';
import styles from "./control.css";

export default function Control({ devices, selectedDevice, deviceSelect, moveCamera, zoomCamera, stepSet, faceDetect, isPtz }) {
  const videoURL = `http://${window.location.hostname}:8888/api/v1/video`;
  const studioURL = `http://${window.location.hostname}:8888/api/v1/studio`;
  console.log("ptz status:", isPtz)

  return (
    <div className="container center-container">
      <div className="inner-content">
        <div className="row">
          <div className="col-md-7">
            <img src={videoURL} alt="Изображение" className="img-fluid video" />
            {/* <img src={studioURL} alt="Изображение" className="img-fluid video" /> */}
          </div>
          <div className="col-md-5">
            <div className="text-center">
              <h3>1234 - Камера конференции</h3>
            </div>
            <hr></hr>
            { isPtz ?
            <div>
              <div className="joystick-container">
                <div>
                  <button className="joystick-button" onClick={() => moveCamera("top")}>↑</button>
                </div>
                <div>
                  <button className="joystick-button" onClick={() => moveCamera("left")}>←</button>
                  <button className="joystick-button" onClick={() => moveCamera("center")}>X</button>
                  <button className="joystick-button" onClick={() => moveCamera("right")}>→</button>
                </div>
                <div>
                  <button className="joystick-button" onClick={() => moveCamera("bottom")}>↓</button>
                </div>
              </div>
              <label for="zoom">Приближение:</label>
              <input type="range" id="zoom" name="zoom" min="1" max="10" className="zoomRange" defaultValue={0} onChange={(event) => zoomCamera(event)} />
              <hr></hr>
            </div>
            : <span style={{color: "#ff0000"}}>Не PTZ камера</span>}
            <h6>Настройки камеры</h6>
            <label className="pr-2">Камера:</label>
            <select onChange={deviceSelect} value={selectedDevice}>
              {Object.entries(devices).map(([key, value]) => (
                <option key={key} value={value}>
                  {key}
                </option>
              ))}
            </select>
            <br></br>
            <label className="pr-2">Шаг:</label>
            <input type="number" min="10" max="500" onChange={(event) => stepSet(event)} defaultValue={"10"}></input>
            <br></br>
            <label className="pr-2">Отслеживание лица:</label>
            <input type="checkbox" onChange={(event) => faceDetect(event)}></input>
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
