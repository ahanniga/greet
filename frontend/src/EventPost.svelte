<script>
    import {EventsEmit} from "../wailsjs/runtime/runtime.js";
    import {GetTaggedProfiles, GetTaggedEvents, GetContactProfile, DeleteEvent} from "../wailsjs/go/main/App.js";
    import { eventStore } from './store.js'
    import LookupPk from "./LookupPk.svelte";
    import LookupEvent from "./LookupEvent.svelte";
    import loadingGif from "./assets/images/loading.gif"

    export let event;

    export let myPk;

    let taggedEvents = [];
    let showWaiting = false;
    let notFound = false;

    const parseContent = (txt) => {
        return imageParse(newlineParse(httpLinkParse(nostrNpubLinkParse(nostrEventLinkParse(txt)))));
    }

    const imageParse = (s) => {
        if(!s) {
            return "";
        }
        return s.replace(/(https?:\/\/\S+(?:jpe?g|png|bmp|gif|webp))/g, '<div><img style="display: inline-block; text-align: left; width: auto; height: auto; max-width: 700px; max-height: 700px; margin: 5px 0; " src="$1" /></div>');
    }

    const newlineParse = (s) => {
        if(!s) {
            return "";
        }
        return s.replace(/\n/g, "<br />");
    }

    const httpLinkParse = (s) => {
        if(!s) {
            return "";
        }
        return s.replace(/https?:\/\/(?!\S+(?:jpe?g|png|bmp|gif|webp))\S+/g, '<a href="#" onclick=runtime.BrowserOpenURL("$&") >$&</a>');
    }

    const nostrEventLinkParse = (s) => {
        if(!s) {
            return "";
        }
        return s.replace(/nostr:(\S[a-zA-Z0-9]+)/g, "<a href='#' data-bs-toggle=\"modal\" data-bs-target=\"#eventDialog\" onclick='runtime.EventsEmit(\"evEventDialog\", \"$1\");'> $& </a>");
    }
    const nostrNpubLinkParse = (s) => {
        if(!s) {
            return "";
        }
        return s.replace(/nostr:(npub\S[a-zA-Z0-9]+)/g, "<a href='#' data-bs-toggle=\"modal\" data-bs-target=\"#profileCard\" onclick='runtime.EventsEmit(\"evProfileCardPk\", \"$1\");'> $& </a>");
    }

    const profileCard = (profile) => {
        EventsEmit("evProfileCard", profile);
    }

    const eventInfo = () => {
        EventsEmit("evEventInfo", event);
    }

    const eventDelete = () => {
        $eventStore = $eventStore.filter(ev => ev.id !== event.id);
        DeleteEvent(event.id);
    }

    const getDisplayName = (profile) => {
        if(!profile || !profile.meta) {
            return "";
        }
        return profile.meta.display_name || profile.meta.name || profile.pk || "";
    }

    const includesMe = (event) => {
        for(let a = 0; event.tags && a < event.tags.length; a++) {
            let tag = event.tags[a];
            if(tag[0] === "p" && tag[1] === myPk) {
                return true;
            }
        }
        return false;
    }

    const getTaggedEventIds = (ev) => {
        let evs = [];
        for(let a = 0; ev.tags && a < ev.tags.length; a++) {
            let tag = ev.tags[a];
            if(tag[0] === "e") {
                evs.push(tag)
            }
        }
        return evs;
    }

    const expandPost = () => {
        if(taggedEvents.length > 0) {
            taggedEvents = [];
            taggedEvents = taggedEvents;
            return;
        }
        showWaiting = true;
        GetTaggedEvents(event.id).then((tes) => {
            for(let a = 0; a < tes.length; a++) {
                let ev = tes[a];
                taggedEvents.push(ev);
            }
            taggedEvents = taggedEvents;
            if(taggedEvents.length === 0) {
                notFound = true;
            }
        }).finally(()=>{
            showWaiting = false;
        });
    }

    const openEventDialog = (id) => {
        EventsEmit("evEventDialog", id);
    }

    const confirmBoost = (ev, by) => {
        EventsEmit("evBoostConfirmDialog", ev, by);
    }

    const openReplyDialog = (ev) => {
        window.runtime.EventsEmit("evReplyDialog", ev);
    }

</script>
<style></style>

<div class="card d-block m-1">
    <!--{@debug event}-->
    {#await GetContactProfile(event.pubkey)}
        <img src="{loadingGif}" width="24" height="24">
    {:then p}
        <div class="card-body p-1 pt-3 pe-0">
        <img src="{p.meta.picture}" alt="" style="width: 36px !important; height: 36px !important; min-width: 36px; min-height: 36px">
        <span class="d-inline ms-2 text-body">{getDisplayName(p)}</span>
        <span class="d-inline me-2 ms-3 float-end small text-body">{event.when}</span>

        <div id="tooltip-container" style="" class="float-end text-muted">
            <a href="#" data-bs-toggle="modal" data-bs-target="#profileCard" data-bs-placement="bottom" title="Contact Info" class="d-inline-block pe-2 nav-link" on:click={() => profileCard(p)}>
                <i class="mb-3 bi bi-person"></i>
            </a>
            <a href="#" data-bs-toggle="modal" data-bs-target="#eventInfo" data-bs-placement="bottom" title="Event Source" class="d-inline-block pe-2 nav-link"  on:click={eventInfo} >
                <i class="mb-3 bi bi-info-circle"></i>
            </a>
            <a href="#" data-bs-toggle="modal" data-bs-target="#replyDialog" data-bs-placement="bottom"  title="Reply" class="d-inline-block pe-2 nav-link" on:click={() => openReplyDialog(event)}>
                <i class="mb-3 bi bi-reply"></i>
            </a>
<!--                <a href="#" data-bs-container="#tooltip-container" data-bs-toggle="tooltip" data-bs-placement="bottom" title="Delete" class="d-inline-block pe-2 nav-link ">-->
<!--                    <i class="mb-3 bi bi-lightning"></i>-->
<!--                </a>-->
<!--                <a href="#" data-bs-toggle="modal" data-bs-placement="bottom" title="Like" class="d-inline-block pe-2 nav-link" on:click={() => { }}>-->
<!--                    <i class="mb-3 bi bi-heart"></i>-->
<!--                </a>-->
            <a href="#" data-bs-toggle="modal" data-bs-target="#confirmDialog" data-bs-placement="bottom" title="Boost" class="d-inline-block pe-2 nav-link" on:click={() => { confirmBoost(event, getDisplayName(p)) }}>
                <i class="mb-3 bi bi-arrow-repeat"></i>
            </a>
            <!--{@debug myPk}-->
            {#if myPk === p.pk}
            <a href="#" data-bs-toggle="modal" data-bs-placement="bottom" title="Delete" class="d-inline-block pe-2 nav-link" on:click={eventDelete}>
                <i class="mb-3 bi bi-trash"></i>
            </a>
            {/if}

        </div>

            <span class="d-inline mb-2 text-muted">{p.meta.nip05 || ""}</span>
        {#if event.kind === 6}
            <span class="badge bg-danger ms-2">Boosted</span>
        {/if}
        {#if includesMe(event) }
            <span class="badge bg-success ms-2">Tagged</span>
        {/if}

        <p class="mt-2 text-primary small">

            {#await GetTaggedProfiles(event.id)}
                <img src="{loadingGif}" width="18" height="18">
            {:then profs}
                {#each profs as prof}
                    <LookupPk {prof} />
                {/each}
            {/await}
            <!--{#each event.tags as tag}-->
            <!--    {#if tag[0] === "t"}-->
            <!--    <span class="text-warning">#{tag[1]} </span>-->
            <!--    {/if}-->
            <!--{/each}-->
        </p>

        {#if event.kind === 1 }
            <p class="font-12 my-3" >
                {@html parseContent(event.content)}
            </p>
        {/if}

        {#if event.kind === 6}
            {#await GetTaggedEvents(event.id)}
                <img src="{loadingGif}" width="18" height="18">
            {:then taggedEvents}
                {#each taggedEvents as taggedEvent}
                <LookupEvent {taggedEvent} {myPk} />
                {/each}
            {/await}
        {:else}
            {#if getTaggedEventIds(event).length > 0}
                <div class="text-muted m-0 p-0 expander" align="left" title="Expand tagged events" on:click={expandPost}>
                    {getTaggedEventIds(event).length}
                    {#if taggedEvents.length === 0}
                        {#if showWaiting}
                    <img style="margin-bottom: 3px" src="{loadingGif}" width="16" height="16">
                        {:else}
                            <i class="bi-caret-right-fill me-2"/>
                        {/if}
                        {#if notFound}
                            <span class="text-danger">Event not found</span>
                        {/if}
                    {:else}
                        <i class="bi-caret-down-fill me-2"/>
                    {/if}
                </div>
            {/if}
        {/if}
        {#each taggedEvents as taggedEvent}
            <LookupEvent {taggedEvent} {myPk} />
        {/each}
        </div>
    {/await}
</div>
