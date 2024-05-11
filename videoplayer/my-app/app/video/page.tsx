
"use strict";
"use client"; // to ensure client side execution as we are using client-side hooks like useRef
import styles from "../page.module.css";
import { useRef } from "react";
import videojs from "video.js";
import VideoJS from "../components/VideoPlayer";
import { useAppContext } from "../page";

export default function Page() {
  let {videoURL,setVideoURL } = useAppContext()
    const playerRef = useRef(null);

    const videoJsOptionsM3u8 = {
      controls: true,
      autoplay: false,
      width: 1500,
      sources: [
        {
          src: videoURL,
          type:  'application/dash+xml'
        },
      ],
      plugins: {
        httpSourceSelector:
        {
          default: 'auto'
        }
      }
    };
    
    const handlePlayerReady = (player: any) => {
      playerRef.current = player;
  
      console.log(player.qualityLevels())
  
      // You can handle player events here, for example:
      player.on('waiting', () => {
        videojs.log('player is waiting');
      });
  
      player.on('dispose', () => {
        videojs.log('player will dispose');
      });
    };
    return (
        <main >
        <div className={styles.container}>
          <VideoJS options={videoJsOptionsM3u8} onReady={handlePlayerReady} />
        </div>
      </main>
    );
}