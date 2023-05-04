<script>
    /**
     *  A small pop-up dialog where users can paste in a hex or npub and view a user profile.
     *  From there that can view recent posts and follow.
     */

    import { Nip19Decode, GetContactProfile } from "../wailsjs/go/main/App.js";
    import { EventsEmit } from "../wailsjs/runtime/runtime.js";

    let input;

    const onFindContactDialog = () => {
        document.getElementById('launchFindContactDialog').click();
        input = "";
        setTimeout(() => {
            document.getElementById('searchContact').focus();
        }, 500);
    }
    window.runtime.EventsOn('evFindContactDialog', onFindContactDialog);

    async function getProfile(pk) {
        return await GetContactProfile(pk);
    }

    const showError = (msg) => {
        let d = document.getElementById("searchErrorMessage");
        d.classList.remove("visually-hidden");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
            d.classList.add("visually-hidden");
        }, 5000);
    }

    async function getNip19Decode(npub) {
        return await Nip19Decode(npub);
    }

    const openProfileCard = (p) => {
        document.getElementById("closeFindContactDialog").click();
        document.getElementById('searchContact').value = "";
        EventsEmit("evProfileCard", p);
    }

    const search = () => {
        let name = document.getElementById('searchContact').value;
        if(!name || name.length < 3) {
            showError("Search too short (minimum 3 characters)");
            return;
        }

        // Npub?
        if(name.startsWith("npub")) {
            console.log("Starts with npub");
            getNip19Decode(name).then((pk) => {
                getProfile(pk).then((p) => {
                    if(p) {
                        openProfileCard(p);
                    }
                });
            }).catch((error) => {
                console.error(error);
                showError("Unable to decode key");
            });
            return;
        }

        // Hex PK?
        if(name.length === 64 && name.indexOf(" ") < 0) {
            getProfile(name).then((p) => {
                if(p) {
                    openProfileCard(p);
                }
            });
            return;
        }

        showError("Not found: " + name);
    }

</script>
<style></style>

<a id="launchFindContactDialog" class="visually-hidden" data-bs-toggle="modal" data-bs-target="#findContactDialog"></a>
<div class="modal fade" id="findContactDialog" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="exampleModalLabel">Find Contact</h5>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">

                <div class="input-group mb-3">
                    <input class="form-control" id="searchContact" placeholder="NIP19 (npub) or Hex key...">
                </div>

            </div>
            <div class="modal-footer">
                <label id="searchErrorMessage" class="me-auto text-danger visually-hidden"></label>
                <button id="closeFindContactDialog" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="button" class="btn btn-primary btn-sm" on:click={search}>Search</button>
            </div>
        </div>
    </div>
</div>