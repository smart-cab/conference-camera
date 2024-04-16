import React, { Component } from "react";

export default function Waiting({ user }) {
    return (
        <div className="container h-100">
          <div className="row h-100 justify-content-center align-items-center mt-5">
            <div className="col-12 text-center">
              <h1>{process.env.REACT_APP_SCHOOL} — Камера конференции</h1>
              <div className="mt-5">
                <h3>Connected: {user}</h3>
              </div>
            </div>
          </div>
        </div>
    );
}