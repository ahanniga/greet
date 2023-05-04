<script>
    /**
     *  A dialog for the creation of a new NOSTR account.
     *  Just the basics are required here, dfurther details can be added in the profile dialog.
     */

    import {GenerateKeys, SaveNewKeys } from "../wailsjs/go/main/App.js";
    import {EventsEmit} from "../wailsjs/runtime/runtime.js";

    let enableSave = true;

    const onGenkeysDialog = () => {
        document.getElementById('launchGenKeysDialog').click();
    }
    window.runtime.EventsOn('evGenkeysDialog', onGenkeysDialog);

    const showError = (msg) => {
        let d = document.getElementById("getKeysErrorMessage");
        d.classList.remove("visually-hidden");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
            d.classList.add("visually-hidden");
        }, 5000);
    }
    const showInfo = (msg) => {
        let d = document.getElementById("getKeysInfoMessage");
        d.classList.remove("visually-hidden");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
            d.classList.add("visually-hidden");
        }, 5000);
    }


    const saveKeys = () => {
        let name = document.getElementById("genkeyName").value.trim();
        let displayName = document.getElementById("genkeyDispName").value.trim();
        let pin = document.getElementById("genKeyPin").value.trim();
        let key = document.getElementById("genPrivKeyInput").value.trim();

        if(key === "") {
            showError("Generate a key before saving!")
            return;
        }
        if(name === "") {
            showError("A name is mandatory")
            return;
        }

        showInfo("Setting...");

        let creds = {
            name: name,
            displayName: displayName,
            pin: pin,
            key: key
        }

        SaveNewKeys(creds).then(()=>{
            document.getElementById("genKeysClose").click();
            EventsEmit("evSuggestFollowsDialog");
        }).catch((msg)=>{
            showError(msg);
        });
    }

    const genKeys = () => {
        GenerateKeys().then((map)=>{
            console.log("Got new keys: PK " + map["pk"]);
            document.getElementById("genPrivKeyInput").value = map["key"];
            document.getElementById("genPubKeyInput").value = map["pk"];
        }).catch((msg)=>{
           showError(msg);
        });
    }

</script>
<style></style>

<a id="launchGenKeysDialog" class="visually-hidden" data-bs-toggle="modal" data-bs-target="#getKeysDialog"></a>
<div class="modal" id="getKeysDialog" tabindex="-1" aria-hidden="true" data-bs-backdrop="static" data-bs-keyboard="false">
    <div class="modal-dialog modal-dialog-centered modal-lg" >
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" ><i class="bi-box-arrow-in-right me-3"></i>New Account</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <legend>About NOSTR Keys</legend>
                Participants in the NOSTR network use keys, not user names and passwords.

                Here you can generate a new public/private key pair. Your account is secured with the private key which should be kept secret.

                <div align="center" class="w-100 mt-4">
                <button type="button" class="btn btn-primary btn-sm ms-3" on:click={genKeys} >Generate Keys</button>
                </div>

                <div class="row">
                    <div class="col">
                        <div class="mb-3">
                            <label for="genkeyName" class="form-label">Name (required)</label>
                            <input type="text" class="form-control" id="genkeyName" placeholder="Random User">
                        </div>
                    </div>
                    <div class="col">
                        <div class="mb-3">
                            <label for="genkeyName" class="form-label"></label>
                            <div id="genkeyName" class="form-text">A name to identify yourself</div>
                        </div>
                    </div>
                </div>
                <div class="row">
                    <div class="col">
                        <div class="mb-3">
                            <label for="genkeyDispName" class="form-label">Display Name (optional)</label>
                            <input type="text" class="form-control" id="genkeyDispName" placeholder="☆¸.•°”˜ ⚡ Random User ˜”°•.¸☆">
                        </div>
                    </div>
                    <div class="col">
                        <div class="mb-3">
                            <label for="genkeyDispNameDesc" class="form-label"></label>
                            <div id="genkeyDispNameDesc" class="form-text">Some clients support display name. This is the best place if you want to use emojis with your name</div>
                        </div>
                    </div>
                </div>

                <div class="mb-2">
                    <label for="genPrivKeyInput" class="form-label">Private Key</label>
                    <input type="text" class="form-control" id="genPrivKeyInput" readonly>
                </div>
                <div class="mb-2">
                    <label for="genPubKeyInput" class="form-label">Public Key</label>
                    <input type="text" class="form-control" id="genPubKeyInput" readonly>
                </div>

                <div class="row">
                    <div class="col">
                        <div class="mb-3">
                            <label for="genKeyPin" class="form-label">PIN (optional)</label>
                            <input type="text" class="form-control" id="genKeyPin">
                        </div>
                    </div>
                    <div class="col">
                        <div class="mb-3">
                            <label for="getKeysDesc" class="form-label"></label>
                            <div id="getKeysDesc" class="form-text">This will be used to encrypt your private key. You will be prompted for the PIN when opening the app.</div>
                        </div>
                    </div>
                </div>

            </div>
            <div class="modal-footer">
                <label id="getKeysErrorMessage" class="ms-lg-2 text-danger visually-hidden"></label>
                <label id="getKeysInfoMessage" class="ms-lg-2 visually-hidden"></label>
                <button id="genKeysClose" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary btn-sm" on:click={saveKeys} {enableSave} >Save</button>
            </div>
        </div>
    </div>
</div>