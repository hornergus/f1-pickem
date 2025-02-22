import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import ThemedApp from "./App";
import reportWebVitals from "./reportWebVitals";
import { CssBaseline } from "@material-ui/core";
import { Provider } from "react-redux";
import store from "./store/store";
import { getLeagues } from "store/actions/leaguesActions";
import { getRaces } from "store/actions/racesActions";

// pre-fetch to optimize page rendering
store.dispatch(getLeagues())
store.dispatch(getRaces((new Date()).getFullYear()))

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <CssBaseline />
      <ThemedApp />
    </Provider>
  </React.StrictMode>,
  document.getElementById("root")
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
