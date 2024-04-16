import React, { useEffect, useRef } from 'react';

function VideoPlayer() {
  const videoRef = useRef(null);

  useEffect(() => {
    const videoElement = videoRef.current;

    return () => {
    };
  }, []);

  return (
    <div>
      <img id="video-websocket" />
    </div>
  );
}

export default VideoPlayer;
