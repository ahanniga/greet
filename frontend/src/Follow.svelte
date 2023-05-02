<script>
    import LookupFollow from "./LookupFollow.svelte";
    import {EventsEmit} from "../wailsjs/runtime/runtime"

    export let profile;

    const profileCard = () => {
        EventsEmit("evProfileCard", profile);
    }

    const getDisplayName = (profile) => {
        return profile.meta.nip05 || "";
    }

    const filter = () => {
        window.runtime.EventsEmit("evFilterByProfile", profile);
    }

    const copyPubkey = () => {
        navigator.clipboard.writeText(profile.pk);
    }
    const copyNpub = () => {
        navigator.clipboard.writeText(profile.npub);
    }

</script>
<style></style>

<div class="d-flex pt-2 pb-2 border-bottom">
    <img src="{profile.meta.picture}" alt="" style="width: 48px !important; height: 48px !important; min-width: 48px; min-height: 48px">

    <p class="mx-lg-2 mb-0 w-100 overflow-x-hidden">
        <LookupFollow { profile } /><br>
        <span class="small" >{getDisplayName(profile)}</span>
    </p>

    <div class="dropdown card-widgets float-end">
        <a href="#" data-bs-toggle="dropdown" aria-expanded="false" class="nav-link">
            <i class="bi bi-three-dots-vertical"></i>
        </a>
        <div class="dropdown-menu dropdown-menu-end">
            <a href="#" class="dropdown-item" data-bs-toggle="modal" data-bs-target="#profileInfo" on:click={profileCard} ><i class="bi bi-person me-2"></i>Profile</a>
            <a href="#" class="dropdown-item" on:click={filter} ><i class="bi bi-filter me-1"></i>Recent Posts</a>
            <li>
                <hr class="dropdown-divider">
            </li>
            <a href="#" class="dropdown-item" on:click={copyPubkey} ><i class="bi bi-clipboard me-2"></i>Copy Pubkey</a>
            <a href="#" class="dropdown-item" on:click={copyNpub} ><i class="bi bi-clipboard me-2"></i>Copy Npub</a>

        </div>
    </div>
</div>



