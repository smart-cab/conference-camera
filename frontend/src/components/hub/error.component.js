import React, { Component } from "react";

export default function Error({ error }) {
    return (
        <div className="container h-100">
          <div className="row h-100 justify-content-center align-items-center mt-5">
            <div className="col-12 text-center">
              <h1>{error}</h1>
            </div>
          </div>
        </div>
    );
}