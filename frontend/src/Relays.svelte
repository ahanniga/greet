<script>
    /**
     *  A dialog to configure relays.
     *  Each can be enabled/disabled and set for read/write.
     *  If set, the current relays are disconnected and new connections are established.
     *  Relay config is written to the greet/config.json file.
     */


    import {GetRelays, SetRelays} from "../wailsjs/go/main/App.js";
    let relays = [];

    const onRelayDialog = () => {
        GetRelays().then((r) => {
            relays = r;
        });
    }
    window.runtime.EventsOn('evRelayDialog', onRelayDialog);

    const showError = (msg) => {
        let d = document.getElementById("relayErrorMessage");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
        }, 5000);
    }
    const showInfo = (msg) => {
        let d = document.getElementById("relayErrorMessage");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
        }, 5000);
    }

    const addRelay = () => {
        let name = document.getElementById('addRelay').value;
        if(!name || name.length < 8) {
            showError("Invalid relay name");
            return;
        }
        name = name.trim().toLowerCase();
        if(!name.match(/^wss?:\/\/\S+/)) {
            showError("Relays start with ws:// or wss://");
            return;
        }
        for(let a = 0; a < relays.length; a++) {
            if(relays[a].url === name) {
                showError("Duplicate name " + name)
                return;
            }
        }
        relays.push({
            url: name,
            enabled: false,
            read: true,
            write: true
        })
        relays = relays;
        document.getElementById('addRelay').value = "";
    }

    const removeRelay = (index) => {
        relays.splice(index, 1);
        relays = relays;
    }
    const relayInfo = (index) => {
        // pop-over? stacked dialog?
    }

    const setRelays = () => {
        for(let a = 0; a < relays.length; a++) {
            let relay = relays[a];
            relay.enabled = document.getElementById("relayEnabled" + a).checked;
            relay.read = document.getElementById("relayRead" + a).checked;
            relay.write = document.getElementById("relayWrite" + a).checked;
        }
        showInfo("Applying new config...")
        SetRelays(relays).then(() => {
            document.getElementById("closeRelaysDialog").click();
        });
    }

</script>
<style></style>

<div class="modal" id="relayDialog" tabindex="-1" aria-hidden="true" data-bs-backdrop="static" data-bs-keyboard="false">
    <div class="modal-dialog modal-dialog-centered modal-lg" >
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" ><i class="bi-hdd-network me-3"></i>Relays</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">

        <form id="relayForm">
            <div class="row">
                <div class="col">
                    <table class="table ">
                        <tbody>
                        {#each relays as relay, i}
                            <tr>
                                <td scope="row">
                                    <div class="form-check form-switch">
                                        <input class="form-check-input" type="checkbox" checked={relay.enabled} id="relayEnabled{i}">
                                        <label class="form-check-label{i}" for="relayEnabled{i}"></label>
                                    </div>
                                </td>
                                <td scope="row" style="max-width: 200px; overflow: hidden">{relay.url}</td>
                                <td scope="row">
                                    <input class="form-check-input" type="checkbox" checked={relay.read} id="relayRead{i}">
                                    <label class="form-check-label" for="relayRead{i}">Read</label>
                                </td>
                                <td scope="row">
                                    <input class="form-check-input" type="checkbox" checked={relay.write} id="relayWrite{i}">
                                    <label class="form-check-label" for="relayWrite{i}">Write</label>
                                </td>
                                <td scope="row"><a href="#" on:click={ ()=>{relayInfo(i)} }><i class="bi-info-circle text-muted"/></a></td>
                                <td scope="row"><a href="#" on:click={ ()=>{removeRelay(i)} } ><i class="bi-trash text-muted"/></a></td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>

                    <div class="input-group mb-3">
                        <input class="form-control" id="addRelay" placeholder="wss://relay.address">
                        <button type="button" class="btn btn-outline-secondary" on:click={addRelay}>Add</button>
                    </div>

                </div>
            </div>
        </form>

            </div>
            <div class="modal-footer">
                <label id="relayErrorMessage" class="ms-3" style="position: absolute; left: 0;"></label>
                <button id="closeRelaysDialog" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary btn-sm" on:click={setRelays}>OK</button>
            </div>
        </div>
    </div>
</div>