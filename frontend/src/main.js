import './assets/bootstrap.bundle.min'
import './assets/bootstrap.min.css'
import './assets/style.css'
import App from './App.svelte'
import {EventsEmit, LogInfo} from "../wailsjs/runtime/runtime.js";

const app = new App({
  target: document.getElementById('app')
})

window.handleNostrLink = function(link) {
  LogInfo("window.handleNostrLink");

  var event = new CustomEvent("onHandleNostrLink", { "detail": link });
  document.dispatchEvent(event);
}


export default app
