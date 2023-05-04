<script>
    /**
     *  Dialog to show details of a user profile.
     *  If profile is of the current user, the fields are editable and
     *  can be saved/published as metadata when saved.
     */

    import {FollowContact, UnfollowContact, GetContactProfile, GetMyPubkey, SaveProfile} from "../wailsjs/go/main/App.js";
    import {EventsEmit} from "../wailsjs/runtime/runtime.js";

    let id = "";
    let npub = "";
    let name = "";
    let picture = "";
    let about = "";
    let nip05 = "";
    let display_name = "";
    let lud06 = "";
    let lud16 = "";
    let website = ""
    let banner = "";
    let following = false;

    let prof;
    let myPk;
    let readonly = true;
    const labelw = 120;
    let changed = false;

    const onProfileCard = (profile) => {
        document.getElementById('launchProfileCardDialog').click();
        GetMyPubkey().then((p)=>{
            myPk = p;
            readonly = p !== profile.pk;
        });

        prof = profile;
        id = profile.pk;
        npub = profile.npub;
        name = profile.meta.name || "";
        display_name = profile.meta.display_name || name;
        picture = profile.meta.picture || "";
        about = profile.meta.about || "";
        nip05 = profile.meta.nip05 || "";
        lud16 = profile.meta.lud16 || "";
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

    const changeCheck = (e) => {
        changed = true;
        name = document.getElementById("prName").value;
        display_name = document.getElementById("prDisplayName").value;
        picture = document.getElementById("prPicture").value;
        banner = document.getElementById("prBanner").value;
        nip05 = document.getElementById("prNip05").value;
    }

    const saveChanges = () => {
        let ln0616 = document.getElementById("prLn").value;
        let hasAt = ln0616.indexOf("@") >= 0;

        lud06 = hasAt ? "" : ln0616;
        lud16 = hasAt ? ln0616 : "";
        about = document.getElementById("prAbout").value;
        website = document.getElementById("prWebsite").value;

        let meta = {
            name: name,
            about: about,
            picture: picture,
            nip05: nip05,
            display_name: display_name,
            lud06: lud06,
            lud16: lud16,
            banner: banner,
            website: website
        }
        SaveProfile(meta).then(()=>{
            changed = false;
        });
    }

    $: lud0616 = lud06.trim() === "" ? lud16 : lud06;

</script>

<style></style>

<a id="launchProfileCardDialog" class="visually-hidden" data-bs-toggle="modal" data-bs-target="#profileCard"></a>
<div class="modal" id="profileCard" tabindex="-1" data-bs-backdrop="static" aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" id="staticBackdropLabel">
                    <div style="height: 30px !important;">
                        <i class="bi-person text-muted me-2"/>
                        <span class="text-muted">{display_name} ({name})</span>
                        <span class="d-inline text-primary" >{nip05}</span>
                    </div>
                </h1>
                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <div class="card-img bg-body-tertiary overflow-y-hidden mb-3" style="height: 250px !important; display: flex; justify-content: center; align-items: center;">
                    <img src="{picture}" alt="" style="position: absolute; top: 30px; left: 30px; width: 120px !important ; height: 120px !important">
                    <img style="width: 100% !important; overflow-y: hidden" src="{banner}">
                </div>

                <div class="input-group mb-3" style="">
                    <span class="input-group-text" style="width: {labelw}px">About</span>
                    <textarea id="prAbout" readonly="{readonly}" class="form-control" on:change={changeCheck}>{about}</textarea>
                </div>

                <div class="input-group mb-3">
                    <span class="input-group-text" style="width: {labelw}px" id="basic-name">Name</span>
                    <input id="prName" readonly="{readonly}" type="text" class="form-control" on:change={changeCheck} value="{name}">
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text" style="width: {labelw}px" id="basic-display-name">Display Name</span>
                    <input id="prDisplayName" readonly="{readonly}" type="text" class="form-control" on:change={changeCheck} value="{display_name}">
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text" style="width: {labelw}px" id="basic-id">Pubkey</span>
                    <input readonly="true" type="text" class="form-control" on:change={changeCheck} value="{id}">
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text" style="width: {labelw}px" id="basic-npub">Npub</span>
                    <input readonly="true" type="text" class="form-control" on:change={changeCheck} value="{npub}">
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text" style="width: {labelw}px" id="basic-picture">Picture</span>
                    <input id="prPicture" readonly="{readonly}" type="text" class="form-control" on:change={changeCheck} value="{picture}">
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text" style="width: {labelw}px" id="basic-banner">Banner</span>
                    <input id="prBanner" readonly="{readonly}" type="text" class="form-control" on:change={changeCheck} value="{banner}">
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text" style="width: {labelw}px" id="basic-website">Website</span>
                    <input id="prWebsite" readonly="{readonly}" type="text" class="form-control" on:change={changeCheck} value="{website}">
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text" style="width: {labelw}px" id="basic-nip05">NIP05</span>
                    <input id="prNip05" readonly="{readonly}" type="text" class="form-control" on:change={changeCheck} value="{nip05}">
                </div>
                <div class="input-group mb-3">
                    <span class="input-group-text" style="width: {labelw}px" id="basic-addon1">LN</span>
                    <input id="prLn" readonly="{readonly}" type="text" class="form-control" on:change={changeCheck} value="{lud0616}">
                </div>
            </div>
            <div class="modal-footer">
                {#if following}
                    <button type="button" class="btn btn-danger btn-sm ms-3" style="position: absolute; left: 0;" on:click={() => { unfollowContact(id) }} >Unfollow</button>
                {:else}
                    <button type="button" class="btn btn-success btn-sm ms-3" style="position: absolute; left: 0;" on:click={() => { followContact(id) }}>Follow</button>
                {/if}
                <button type="button" class="btn btn-success btn-sm" data-bs-dismiss="modal" on:click={filter}>Recent Posts</button>
                {#if !readonly}
                <button type="button" disabled="{!changed}" class="btn btn-success btn-sm" on:click={saveChanges}>Save Changes</button>
                {/if}
                <button id="closeProfileCardDialog" type="button" class="btn btn-primary btn-sm" data-bs-dismiss="modal">Close</button>
            </div>
        </div>
    </div>
</div>