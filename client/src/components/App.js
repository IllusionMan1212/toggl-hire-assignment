import { Provider } from "urql";
import { client } from "../client";
import { Questions } from "./Questions";
import { ToastWrapper } from "../contexts/toastContext";
import { Routes, Route } from "react-router-dom";
import Results from "./Results";

export function App() {
  return (
    <Provider value={client}>
      <ToastWrapper>
        <Routes>
          <Route path="/" element={<Questions/>}/>
          <Route path="/results" element={<Results/>}/>
        </Routes>
      </ToastWrapper>
    </Provider>
  );
}
