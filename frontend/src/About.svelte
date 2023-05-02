<script>
    import nostrIcon from "./assets/images/nostr.png";
    import {BrowserOpenURL} from "../wailsjs/runtime/runtime.js";
</script>
<style></style>

<!-- Modal -->
<div class="modal fade" id="aboutDialog" >
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="exampleModalLabel">About Greet</h5>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body w-100" style="text-align: center;">
                <h4 class="">Greet Nostr Client</h4>
                <img src={nostrIcon} alt="" height="28" class="m-lg-3">
                <div class="">Version 0.1.0</div>
                <div class=""><a on:click={()=>{BrowserOpenURL("https://github.com/ahanniga/greet")}} href="#">https://github.com/ahanniga/greet</a></div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-primary btn-sm" data-bs-dismiss="modal">Close</button>
            </div>
        </div>
    </div>
</div>
