import { Provider } from "urql";
import { client } from "../client";
import { Questions } from "./Questions";
import { ToastWrapper } from "../contexts/toastContext";

export function App() {
  return (
    <Provider value={client}>
      <ToastWrapper>
        <Questions />
      </ToastWrapper>
    </Provider>
  );
}
