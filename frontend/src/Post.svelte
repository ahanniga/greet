<script>
    /**
     *  A modal dialog to create a new note for publishing.
     *  User can add p-tags by selecting existing, or add new
     *  recipients not already followed.
     *
     *  User can also select which relays to broadcast the message.
     */

    import {
        GetWritableRelays,
        Nip19Decode,
        PublishContentToSelectedRelays,
        GetContactProfile
    } from "../wailsjs/go/main/App.js";
    import {EventsEmit} from "../wailsjs/runtime/runtime.js";
    import {contactStore} from "./store.js";

    let promise;
    let eventTags = [];

    const onPostDialog = () => {
        promise = GetWritableRelays();
        document.getElementById('postForm').focus();
    }
    window.runtime.EventsOn('evPostDialog', onPostDialog);

    const showError = (msg) => {
        let d = document.getElementById("postErrorMessage");
        d.classList.remove("visually-hidden");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
            d.classList.add("visually-hidden");
        }, 2000);
    }

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
        let content = document.getElementById('postForm').value;

        for(let a = 0;; a++) {
            let cb = document.getElementById('relayCheck' + a);
            if(cb && cb.checked === true) {
                relays.push(cb.labels[0].textContent);
            }
            else break;
        }

        if(relays.length === 0) {
            showError("Select at least one relay");
            return;
        }
        if(content.length === 0) {
            showError("The content is empty?");
            return;
        }

        PublishContentToSelectedRelays(1, content, eventTags, relays).then(() => {
            EventsEmit("evReloadSavedEvents");
        });
        document.getElementById("closePostDialog").click();
    }

    async function getNip19Decode(npub) {
        return await Nip19Decode(npub);
    }

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
            getNip19Decode(name).then((pk) => {
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

<div class="modal" id="postDialog" tabindex="-1" data-bs-backdrop="static" data-bs-keyboard="false" aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-xl">
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" id="staticBackdropLabel"><i class="bi-chat-text-fill"></i> Post New Note</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">

                <form id="postMessageForm">
                    <!-- Message input -->
                    <div class="form-outline mb-4">
                        <label class="form-label" for="postForm">Message</label>
                        <textarea class="form-control" id="postForm" rows="4" placeholder=""></textarea>
                    </div>

                    <div class="row">
                        <div class="col">
                            <label class="form-label" for="postForm">Tagged People</label>
                            <table class="table ">
                                <tbody>

                                {#each eventTags as et}
                                    {#await GetContactProfile(et[1]) then p}
                                        <!--{@debug p}-->
                                        <tr>
                                            <td scope="row"><img src="{p.meta.picture}" alt="" style="width: 32px !important; min-width: 36px; min-height: 36px"></td>
                                            <td scope="row">{p.meta.display_name || p.meta.name || ""}</td>
                                            <td>{p.meta.nip05 || ""}</td>
                                            <td><a href="#" on:click={()=>{removeTagged(p.pk)}}><i class="bi-trash text-muted"/></a></td>
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
                <label id="postErrorMessage" class="me-auto text-danger visually-hidden"></label>
                <button id="closePostDialog" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary btn-sm" on:click={postMessage}>Post</button>
            </div>
        </div>
    </div>
</div>
