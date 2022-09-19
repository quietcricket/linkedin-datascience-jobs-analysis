(function () {
	document.querySelector('.jobs-search__results-list li:last-child').scrollIntoView({ behavior: "smooth" });
	setTimeout(() => {
		window.scrollBy(0, -Math.random() * 500);
	}, Math.random() * 1000 + 1000);
})();