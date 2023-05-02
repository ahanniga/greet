<script>
    import {Nip19Decode, SetLoginWithPrivKey} from "../wailsjs/go/main/App.js";
    import {EventsEmit} from "../wailsjs/runtime/runtime.js";

    const showError = (msg) => {
        let d = document.getElementById("loginErrorMessage");
        d.classList.remove("visually-hidden");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
            d.classList.add("visually-hidden");
        }, 5000);
    }
    const showInfo = (msg) => {
        let d = document.getElementById("loginInfoMessage");
        d.classList.remove("visually-hidden");
        d.innerText = msg;
        setTimeout(() => {
            d.innerText = "";
            d.classList.add("visually-hidden");
        }, 5000);
    }

    const onLaunchNewAccount = () => {
        EventsEmit("evGenkeysDialog");
    }

    const setLogin = () => {
        let privKeyInput = document.getElementById("privKeyInput").value;
        let pinInput = document.getElementById("pinDialogInput").value;

        privKeyInput = privKeyInput.trim();
        if(privKeyInput === "") {
            showError("Invalid private key")
        }
        if(privKeyInput.startsWith("nsec")) {
            if(privKeyInput.length !== 63) {
                showError("Bad key length")
            }
            Nip19Decode(privKeyInput).then((hex)=> {
                console.log(hex);
            }).catch((e) => {
                console.error(e);
                showError("Unable to decode key");
            });
        } else {
            if(privKeyInput.length !== 64) {
                showError("Bad key length")
            }
        }

        // Looks good
        showInfo("Setting login...");
        SetLoginWithPrivKey([privKeyInput, pinInput]).then(() => {
            document.getElementById("closeLoginDialog").click();
        }).catch((e) => {
            console.error(e);
            showError(e);
        });
    }

</script>
<style></style>

<div class="modal" id="loginDialog" tabindex="-1" aria-hidden="true" data-bs-backdrop="static" data-bs-keyboard="false">
    <div class="modal-dialog modal-dialog-centered modal-lg" >
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" ><i class="bi-box-arrow-in-right me-3"></i>NOSTR Login</h1>
                <button type="button" class="btn-close btn-sm" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                <legend>Welcome to Greet!</legend>
                If you've used a NOSTR client before then enter your private key below. If you're new, then click on the Create Account button to get set up.
                <div class="mb-3 mt-4">
                    <label for="privKeyInput" class="form-label">Private Key</label>
                    <input type="text" class="form-control" id="privKeyInput" placeholder="NIP19 (nsec) or Hex key...">
                </div>

                <div class="row">
                    <div class="col">
                        <div class="mb-3">
                            <label for="pinDialogInput" class="form-label">PIN (optional)</label>
                            <input type="text" class="form-control" id="pinDialogInput">
                        </div>
                    </div>
                    <div class="col">
                        <div class="mb-3">
                            <label for="loginDesc" class="form-label"></label>
                            <div id="loginDesc" class="form-text">This will be used to encrypt your private key. You will be prompted for the PIN when opening the app.</div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <label id="loginErrorMessage" class="ms-lg-2 text-danger visually-hidden"></label>
                <label id="loginInfoMessage" class="ms-lg-2 visually-hidden"></label>
                <button type="button" class="btn btn-success btn-sm ms-3" style="position: absolute; left: 0;" data-bs-dismiss="modal" on:click={onLaunchNewAccount} >Create Account</button>
                <button id="closeLoginDialog" type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="submit" class="btn btn-primary btn-sm" on:click={setLogin}>Login</button>
            </div>
        </div>
    </div>
</div>