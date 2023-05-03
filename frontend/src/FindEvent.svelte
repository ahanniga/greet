<script>
    import { Nip19Decode } from "../wailsjs/go/main/App.js";
    import { EventsEmit } from "../wailsjs/runtime/runtime.js";

    const onFindEventDialog = () => {
        setTimeout(() => {
            document.getElementById('findEvent').focus();
        }, 500);
    }
    window.runtime.EventsOn('evFindEventDialog', onFindEventDialog);


    const showError = (msg) => {
        let d = document.getElementById("findEventErrorMessage");
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

    const close = () => {
        document.getElementById("closeFindEventDialog").click();
    }

    const search = () => {
        let id = document.getElementById('findEvent').value;
        if(!id || id.length < 63 || id.length > 64)  {
            showError("The ID should be 63 or 64 characters long");
            return;
        }

        // Nip19 encoded?
        if(id.startsWith("note")) {
            getNip19Decode(id).then((hex) => {
                EventsEmit("evEventDialog", hex);
                close();
            }).catch((error) => {
                console.error(error);
                showError("Unable to decode the event ID");
            });
        } else {
            EventsEmit("evEventDialog", id);
            close();
        }
    }

</script>
<style></style>

<!-- Modal -->
<div class="modal fade" id="findEventDialog" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">Find Event by ID</h5>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">

                <div class="input-group mb-3">
                    <input class="form-control" id="findEvent" placeholder="NIP19 (note) or Hex ID  ...">
                </div>

            </div>
            <div class="modal-footer">
                <label id="findEventErrorMessage" class="me-auto text-danger visually-hidden"></label>
                <button id="closeFindEventDialog" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="button" class="btn btn-primary btn-sm" on:click={search}>Search</button>
            </div>
        </div>
    </div>
</div>