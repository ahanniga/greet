<script>
    import {FollowContact, UnfollowContact, GetContactProfile} from "../wailsjs/go/main/App.js";
    import {EventsEmit} from "../wailsjs/runtime/runtime.js";

    let id = "";
    let npub = "";
    let name = "";
    let picture = "";
    let about = "";
    let nip05 = "";
    let display_name = "";
    let lud06 = "";
    let website = ""
    let banner = "";
    let following = false;

    let prof;

    const onProfileCard = (profile) => {
        prof = profile;
        id = profile.pk;
        npub = profile.npub;
        name = profile.meta.name || "";
        display_name = profile.meta.display_name || name;
        picture = profile.meta.picture || "";
        about = profile.meta.about || "";
        nip05 = profile.meta.nip05 || "";
        lud06 = profile.meta.lud06 || "";
        website = profile.meta.website || "";
        banner = profile.meta.banner || "";
        following = profile.following;
    }
    window.runtime.EventsOn('evProfileCard', onProfileCard);

    const onProfilePk = (pk) => {
        GetContactProfile(pk).then((p)=>{
            EventsEmit("evProfileCard", p);
        })
    }
    window.runtime.EventsOn('evProfileCardPk', onProfilePk);

    const followContact = (pk) => {
        FollowContact([pk]).then((err) => {
            document.getElementById("closeProfileCardDialog").click();
        });

    }
    const unfollowContact = (pk) => {
        UnfollowContact(pk).then((err) => {
            document.getElementById("closeProfileCardDialog").click();
        });
    }

    const filter = () => {
        window.runtime.EventsEmit("evFilterByProfile", prof);
    }

</script>

<style></style>

<div class="modal" id="profileInfo" tabindex="-1" data-bs-backdrop="static" aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" id="staticBackdropLabel"><i class="bi-person text-muted me-2"/>{display_name}</h1>
                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <img src="{picture}" alt="" style="width: 100px !important ; height: 100px !important">
                <span class="text-muted ms-3">{display_name} {name}</span>
                <span class="d-inline text-primary" >{nip05}</span>
                <p class="mt-2 d-block">
                    <span class="d-inline">{about}</span>
                </p>
                <p class="mt-1 d-block">
                    <span class="d-inline-block w-label-70">Pubkey:</span>
                    <span class="d-inline"><code>{id}</code></span>
                </p>
                <p class="mt-1 d-block">
                    <span class="d-inline-block w-label-70">Npub:</span>
                    <span class="d-inline"><code>{npub}</code></span>
                </p>
                {#if website !== ""}
                <p class="mt-1 d-block">
                    <span class="d-inline-block w-label-70">Website:</span>
                    <span class="d-inline"><code>{website}</code></span>
                </p>
                {/if}
                {#if lud06 !== ""}
                    <p class="mt-1 d-block">
                        <span class="d-inline-block w-label-70">LN:</span>
                        <span class="d-inline"><code>{lud06}</code></span>
                    </p>
                {/if}
            </div>
            <div class="modal-footer">
                {#if following}
                    <button type="button" class="btn btn-danger btn-sm ms-3" style="position: absolute; left: 0;" on:click={() => { unfollowContact(id) }} >Unfollow</button>
                {:else}
                    <button type="button" class="btn btn-success btn-sm ms-3" style="position: absolute; left: 0;" on:click={() => { followContact(id) }}>Follow</button>
                {/if}
                <button type="button" class="btn btn-success btn-sm" data-bs-dismiss="modal" on:click={filter}>Recent Posts</button>
                <button id="closeProfileCardDialog" type="button" class="btn btn-primary btn-sm" data-bs-dismiss="modal">Close</button>
            </div>
        </div>
    </div>
</div>