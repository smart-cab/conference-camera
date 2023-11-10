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
      <img src="http://192.168.1.13:8888/api/v1/video" />
    </div>
  );
}

export default VideoPlayer;
