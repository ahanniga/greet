<script>
    import {
        GetWritableRelays,
        PublishContentToSelectedRelays,
        Nip19Decode,
        GetContactProfile
    } from "../wailsjs/go/main/App.js";
    import { EventsEmit } from "../wailsjs/runtime/runtime.js";
    import { contactStore } from './store.js'

    let promise;
    let event;
    let eventTags = [];

    const showError = (msg) => {
        let d = document.getElementById("replyErrorMessage");
        d.classList.remove("visually-hidden");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
            d.classList.add("visually-hidden");
        }, 2000);
    }

    const onReplyDialog = (ev) => {
        event = ev;
        promise = GetWritableRelays();
        eventTags = [];

        eventTags.push(["p", event.pubkey]);
        for(let a = 0; a < event.tags.length; a++) {
            let tag = event.tags[a];
            if(tag[0] === "p" && !hasExistingPTag(eventTags, tag[1])){
                eventTags.push(tag);
            }
        }
        document.getElementById('replyForm').focus();
    }
    window.runtime.EventsOn('evReplyDialog', onReplyDialog);

    const hasExistingPTag = (tags, val) => {
        for(let a = 0; a < tags.length; a++) {
            let tag = tags[a];
            if(tag[0] === "p") {
                if(tag[1] === val) {
                    return true;
                }
            }
        }
        return false;
    }

    const postMessage = () => {
        let relays = [];
        let content = document.getElementById('replyForm').value;

        for(let a = 0;; a++) {
            let cb = document.getElementById('relayCheck' + a);
            if(cb && cb.checked === true) {
                relays.push(cb.labels[0].textContent);
            }
            else break;
        }
        eventTags.push(["e", event.id]);

        if(relays.length === 0) {
            showError("Select at least one relay");
            return;
        }

        PublishContentToSelectedRelays(1, content, eventTags, relays).then(() => {
            EventsEmit("evReloadSavedEvents");
        });
    }

    // async function getNip19Decode(npub) {
    //     return await Nip19Decode(npub);
    // }

    const removeTagged = (pk) => {
        let tmpTags = [];
        for(let a = 0; a < eventTags.length; a++) {
            let tag = eventTags[a];
            if(tag[0] === "p" && tag[1] !== pk) {
                tmpTags.push(tag)
            }
        }
        eventTags = tmpTags;
    }

    const addTagged = () => {
        let name = document.getElementById('taggedContact').value;
        if(!name || name.length < 3) {
            return;
        }

        // Npub?
        if(name.startsWith("npub")) {
            console.log("Starts with npub");
            Nip19Decode(name).then((pk) => {
                if(!hasExistingPTag(eventTags, pk)) {
                    eventTags.push(["p", pk]);
                    eventTags = eventTags;
                }
                document.getElementById('taggedContact').value = "";
            });
            return;
        }

        // Hex PK?
        if(name.length === 64 && name.indexOf(" ") < 0) {
            if(!hasExistingPTag(eventTags, name)) {
                eventTags.push(["p", name]);
                eventTags = eventTags;
            }
            document.getElementById('taggedContact').value = "";
            return;
        }

        for(let a = 0; a < $contactStore.length; a++) {
            let c = $contactStore[a];
            if(name === c.meta.display_name || name === c.meta.name || name === c.meta.nip05) {
                if(!hasExistingPTag(eventTags, c.pk)) {
                    eventTags.push(["p", c.pk]);
                    eventTags = eventTags;
                }
                document.getElementById('taggedContact').value = "";
                return;
            }
        }
        showError("Not found: " + name);
    }

</script>
<style></style>

<div class="modal" id="replyDialog" tabindex="-1" data-bs-backdrop="static" data-bs-keyboard="false" aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-xl">
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" id="staticBackdropLabel"><i class="bi-reply"></i> Post Reply</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">

                <form id="replyMessageForm">
                    <!-- Message input -->
                    <div class="form-outline mb-4">
                        <label class="form-label" for="replyForm">Message</label>
                        <textarea class="form-control" id="replyForm" rows="4" placeholder=""></textarea>
                    </div>
                    <!--{@debug event}-->
                    <div class="row">
                        <div class="col">
                            <label class="form-label" for="replyForm">Tagged People</label>
                            <table class="table ">
                                <tbody>

                                {#each eventTags as et}
                                    {#await GetContactProfile(et[1]) then p}
                                        <!--{@debug p}-->
                                        <tr>
                                            <td scope="row"><img src="{p.meta.picture}" alt="" style="width: 32px !important; min-width: 36px; min-height: 36px"></td>
                                            <td scope="row">{p.meta.display_name || p.meta.name || ""}</td>
                                            <td>{p.meta.nip05 || ""}</td>
                                            {#if p.pk !== event.pubkey}
                                                <td><a href="#" on:click={()=>{removeTagged(p.pk)}}><i class="bi-trash text-muted"/></a></td>
                                            {:else }
                                                <td>&nbsp;</td>
                                            {/if}
                                        </tr>
                                    {/await}
                                {/each}

                                </tbody>
                            </table>

<!--                            <label for="exampleDataList" class="form-label">Add</label>-->
                            <div class="input-group mb-3">
                                <input class="form-control" list="datalistOptions" id="taggedContact" placeholder="Follows, npub or pubkey...">
                                <datalist id="datalistOptions">
                                    {#each $contactStore as contact}
                                        <!--{@debug contact}-->
                                        <option pk="{contact.pk}" value="{contact.meta.display_name || contact.meta.name || contact.meta.nip05}">
                                    {/each}
                                </datalist>
                                <button type="button" class="btn btn-outline-secondary" on:click={addTagged}>Add</button>
                            </div>

                        </div>
                        <div class="col">
                            <label class="form-label" for="replyForm">Relays</label>

                            {#if promise}
                                {#await promise}
                                {:then relays}
                                    {#if relays.length === 0}
                                        <div align="center" class="mx-2"><i>(None)</i></div>
                                    {/if}
                                    {#each relays as relay, i}
                                        <!--{@debug relay}-->
                                        <div class="form-check form-switch">
                                            <input class="form-check-input" type="checkbox" checked=true id="relayCheck{i}">
                                            <label class="form-check-label" for="relayCheck{i}">
                                                {relay}
                                            </label>
                                        </div>
                                    {/each}
                                {/await}
                            {/if}
                        </div>
                    </div>


                </form>
            </div>
            <div class="modal-footer">
                <label id="replyErrorMessage" class="me-auto text-danger visually-hidden"></label>
                <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary btn-sm" data-bs-dismiss="modal" on:click={postMessage}>Reply</button>
            </div>
        </div>
    </div>
</div>
