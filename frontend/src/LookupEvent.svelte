<script>
    import EventPost from "./EventPost.svelte";
    export let taggedEvent;
    export let myPk;

    $: event = taggedEvent;
</script>
<style></style>

<EventPost {event} {myPk} />

