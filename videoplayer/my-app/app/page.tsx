"use strict";
"use client"; // to ensure client side execution as we are using client-side hooks like useRef
import styles from "./page.module.css";
import { useRef, useState } from "react";
import videojs from "video.js";
import VideoJS from "./components/VideoPlayer";
import { redirect } from "next/navigation";
import { createContext, useContext } from 'react';
import { useRouter } from "next/navigation";

const AppContext = createContext<any>(undefined);

export function AppWrapper({children} : {
  children: React.ReactNode
}){
  const [videoURL, setVideoURL] = useState('');

  return (
    <AppContext.Provider value={{videoURL, setVideoURL} }>
      {children}
    </AppContext.Provider>
  );
}

export function useAppContext (){
  return useContext(AppContext)
}

export default function Home() {
  let {videoURL,setVideoURL } = useAppContext()
  const router = useRouter()

  const handleSubmit = async(e: any) => {
    e.preventDefault()
    router.push("/video")
  }

  return (

    <main >
     <form  onSubmit={handleSubmit} className="w-1/2">
  
        <label>
          <input 
          required
          type="text"
          onChange={(e) => setVideoURL(e.target.value)}
          value={videoURL}
          >
          </input>
        </label>
        <button
        className="btn-primary"
        
        >
            Open Video Player
        </button>

     </form>
    </main>
  );
}