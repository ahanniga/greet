<script>
    import {PostEvent} from "../wailsjs/go/main/App.js";

    let event; // = defaultEvent;
    let name = "unknown"

    const onBoostConfirmDialog = (ev, by) => {
        event = ev;
        name = by;
    }
    window.runtime.EventsOn('evBoostConfirmDialog', onBoostConfirmDialog);

    const postBoost = () => {
        let tags = [];
        tags.push(["e", event.id]);
        tags.push(["p", event.pubkey]);
        PostEvent(6, tags, "").then(() => {
            document.getElementById("boostClose").click();
        });
    }
</script>
<style></style>

<!-- Modal -->
<div class="modal fade" id="confirmDialog" tabindex="-1" data-bs-backdrop="static" aria-labelledby="exampleModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="exampleModalLabel">Confirm</h5>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <p>Boost this event by {name} ?</p>
            </div>
            <div class="modal-footer">
                <button id="boostClose" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="button" class="btn btn-primary btn-sm" on:click={postBoost}>Boost!</button>
            </div>
        </div>
    </div>
</div>