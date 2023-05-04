<script>
    /**
     *  Loaded by Follow, formats the followed contact and launches a ProfileCard when clicked
     */

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
    .profile-link {
        margin-right: 8px;
        text-decoration: none;
        font-size: medium;
        text-transform: capitalize
    }
</style>

<a href="#" class="profile-link" data-bs-toggle="modal" data-bs-target="#profileCard" on:click={profileCard} >{ getDisplayName(profile) }</a>


