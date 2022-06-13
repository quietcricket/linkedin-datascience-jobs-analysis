


(function () {
	document.querySelector('.jobs-search__results-list li:last-child').scrollIntoView({ behavior: "smooth" });
	setTimeout(() => {
		window.scrollBy(0, -20 + Math.random() * 20);
	}, Math.random() * 1000 + 100);
})();