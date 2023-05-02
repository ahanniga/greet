<script>
    import {EventsEmit} from "../wailsjs/runtime/runtime"

    export let profile;
    const profileCard = () => {
        EventsEmit("evProfileCard", profile);
    }
    const getDisplayName = (profile) => {
        return (profile.meta.display_name || profile.meta.name || profile.meta.nip05 || profile.pk || "");
    }

</script>
<style>
    a {
        margin-right: 8px;
        text-decoration: none;
        font-size: medium;
    }
</style>

<a href="#" data-bs-toggle="modal" data-bs-target="#profileInfo" style="text-transform: capitalize;" on:click={profileCard} >{ getDisplayName(profile) }</a>


