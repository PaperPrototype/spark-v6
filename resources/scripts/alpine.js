document.addEventListener("alpine:init", function(event) {
	Alpine.data('dropdown', () => ({
		open: false,
	 
		toggle() {
			this.open = ! this.open
		}
	}));
});