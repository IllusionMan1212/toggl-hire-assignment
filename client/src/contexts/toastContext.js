import { createContext, useContext, useState } from "react";
import Toast from "../components/Toast";

const ToastContext = createContext(null);

export function ToastWrapper({ children }) {
  const [toasts, setToasts] = useState([]);

  const setToastAttributes = (text, duration) => {
    const id = (Date.now().toString(36) + Math.random().toString(36).substr(2, 5)).toUpperCase();

    setToasts((toasts) => {
      return toasts.concat({ text, id });
    });

    setTimeout(() => {
      setToasts((toasts) => {
        return toasts.filter((toasts) => toasts.id !== id);
      });
    }, duration);
  };

  return (
    <>
      <div className="toasts">
        {toasts.map((toast) => {
          return <Toast key={toast.id} text={toast.text}/>
        })}
      </div>
      <ToastContext.Provider value={setToastAttributes}>
        {children}
      </ToastContext.Provider>
    </>
  )
}

export function useToastContext() {
  return useContext(ToastContext);
}
