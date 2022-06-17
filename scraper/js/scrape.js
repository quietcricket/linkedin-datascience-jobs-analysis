
(function () {
	let nodes = document.querySelectorAll('.jobs-search__results-list .base-card__full-link');
	// skip those already scraped
	let result = "";
	for (let i = window.LAST_INDEX ? window.LAST_INDEX : 0; i < nodes.length; i++) {
		let n = nodes[i];
		result += n.textContent.trim() + "," + n.getAttribute('href') + "\n";
	}
	if (result.length > 0) {
		window.LAST_INDEX = nodes.length;
		window.LAST_TIMESTAMP = new Date().getTime();
	} else {
		if (new Date().getTime() - window.LAST_TIMESTAMP > 30000) {
			window.LAST_TIMESTAMP = new Date().getTime();
			result = "stop";
		}
	}
	return result;
})();