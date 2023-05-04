<script>
    /**
     *  Called by EventPost to create clickable links for each p-tag in a note
     */

    import {EventsEmit} from "../wailsjs/runtime/runtime.js";

    export let prof;

    const profileCard = () => {
        EventsEmit("evProfileCard", prof);
    }

</script>
<style>
    a {
        margin-right: 8px;
        text-decoration: none;
    }
</style>

<a href="#" data-bs-toggle="modal" data-bs-target="#profileCard" on:click={profileCard} >{prof.meta.nip05 || prof.pk}</a>&nbsp;



