import {
	createSSRApp
} from "vue";
import App from "./App.vue";
// Import router to trigger initialization (router initializes itself in constructor)
import "./utils/router.js";

export function createApp() {
	const app = createSSRApp(App);

	return {
		app,
	};
}
