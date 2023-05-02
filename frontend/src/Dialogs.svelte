<!--This is a workaround to some strange behaviour between Wails and Bootstrap -->
<!-- https://github.com/wailsapp/wails/discussions/2636 -->

<script>
    import PinDialog from "./PinDialog.svelte";
    import Login from "./Login.svelte";
    import FindContact from "./FindContact.svelte";
    import GenKeys from "./GenKeys.svelte";
    import Post from "./Post.svelte";
    import ProfileCard from "./ProfileCard.svelte";
    import Relays from "./Relays.svelte";
    import SuggestFollows from "./SuggestFollows.svelte";
    import EventInfo from "./EventInfo.svelte";
    import EventDialog from "./EventDialog.svelte";
    import Reply from "./Reply.svelte";
    import ConfirmBoost from "./ConfirmBoost.svelte";
    import About from "./About.svelte";

    const onPinDialog = () => {
        document.getElementById('launchPinDialog').click();
    }
    window.runtime.EventsOn('evPinDialog', onPinDialog);

    const onLoginDialog = () => {
        document.getElementById('launchLoginDialog').click();
    }
    window.runtime.EventsOn('evLoginDialog', onLoginDialog);

    const onFindContactDialog = () => {
        document.getElementById('launchFindContactDialog').click();
    }
    window.runtime.EventsOn('evFindContactDialog', onFindContactDialog);

    const onGenkeysDialog = () => {
        document.getElementById('launchGenKeysDialog').click();
    }
    window.runtime.EventsOn('evGenkeysDialog', onGenkeysDialog);

    const onSuggestFollowsDialog = () => {
        document.getElementById('launchSuggestFollowsDialog').click();
    }
    window.runtime.EventsOn('evSuggestFollowsDialog', onSuggestFollowsDialog);

</script>


<a id="launchPinDialog" data-bs-toggle="modal" data-bs-target="#pinDialog"></a>
<a id="launchLoginDialog" data-bs-toggle="modal" data-bs-target="#loginDialog"></a>
<a id="launchFindContactDialog" data-bs-toggle="modal" data-bs-target="#findContactDialog"></a>
<a id="launchGenKeysDialog" data-bs-toggle="modal" data-bs-target="#getKeysDialog"></a>
<a id="launchSuggestFollowsDialog" data-bs-toggle="modal" data-bs-target="#suggestFollowsDialog"></a>
<PinDialog />
<Login />
<FindContact />
<GenKeys />
<Post />
<ProfileCard />
<Relays />
<SuggestFollows />
<EventInfo />
<EventDialog />
<Reply />
<ConfirmBoost />
<About />




