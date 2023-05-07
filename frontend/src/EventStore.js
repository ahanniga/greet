import {writable, derived } from "svelte/store";
import orderBy from 'lodash/orderBy'

const events_ = writable([]);

export const eventStore = {
    subscribe: events_.subscribe,
    addEvent: ev =>
        events_.update(evs => [
            ev, ...evs
        ]),
    updateEvent: (id, ev) =>
        events_.update(evs =>
            evs.map(ev => (ev.id === id ? { ev, ...evs } : ev))
        ),
    deleteEvent: id =>
        events_.update(evs => evs.filter(ev => ev.id !== id)),
    deleteAll: () => events_.set([]),
};

export const sortedEvents = derived(events_, (events_) => orderBy(events_, ['created_at'], ['desc']))