<script>
    import {GetMyPubkey, GetTextNotesByEventIds, Nip19Decode} from "../wailsjs/go/main/App.js";
    import EventPost from "./EventPost.svelte";
    import defaultEvent from './Util.svelte'

    let myPk;
    let event = defaultEvent;

    const onEventDialog = (noteRef) => {
        event = false;
        GetMyPubkey().then((pk) => {
            myPk = pk;
            if(noteRef.startsWith("note" || noteRef.startsWith("nevent"))) {
                Nip19Decode(noteRef).then((hexId)=>{
                   GetTextNotesByEventIds([hexId]).then((events)=>{
                       if(events.length > 0) {
                           event = events[0];
                       }
                   });
                }).catch((err)=>{
                    console.log("Error:" + err);
                });
            } else {
                GetTextNotesByEventIds([noteRef]).then((events)=>{
                    if(events.length > 0) {
                        event = events[0];
                    }
                });
            }
        });
    }

    window.runtime.EventsOn('evEventDialog', onEventDialog);

</script>
<style></style>

<div class="modal " id="eventDialog" tabindex="-1" data-bs-backdrop="static" aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-xl">
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" id="staticBackdropLabel">Event {event.id}</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                {#if event}
                    <EventPost {event} {myPk}/>
                {/if}
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-primary btn-sm" data-bs-dismiss="modal">Close</button>
            </div>
        </div>
    </div>
</div>
