<script>
    /**
     *  This component is called from an EventPost and simply loads another
     *  (e-tagged) EventPost with an "event" property
     */

    import EventPost from "./EventPost.svelte";
    export let taggedEvent;
    export let myPk;

    $: event = taggedEvent;
</script>
<style></style>

<EventPost {event} {myPk} />

