<script lang="ts">
	let { class: className = 'h-8 w-8', title = 'Standblatt' }: { class?: string; title?: string } = $props();
</script>

<svg class={className} viewBox="0 -75 62 78" fill="currentColor" role="img" aria-label={title} xmlns="http://www.w3.org/2000/svg">
	<circle cx="28" cy="-27" r="23" fill="none" stroke="currentColor" stroke-width="12" />
	<circle cx="28" cy="-27" r="8" />
	<rect x="50" y="-73" width="11" height="73" />
</svg>
