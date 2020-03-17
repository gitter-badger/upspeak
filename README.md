# Upspeak

A real-time, distributed, and contextual discussion platform built on Matrix.

Upspeak aims to be a natural discussion and knowledge sharing system that enables discussions based on context (or, maintains context of all conversations), seamlessly blends synchronous and asynchronous communication modes, and provides a distraction-free experience.

## Project status

Still an idea under experiment.

## Focus areas

### Maintain context of all discussions

Create and retain relations between discussions. Given any point in a conversation, show how it is connected to the rest of the system. Make it easy to access those relations.

### Context graph

Group conversations in threads, rooms and workspaces. Make every entity discoverable and linkable. Entities can be cross-linked with a simple syntax. Visualize the connectedness of conversation. Switch between thread view, graph view, and chat scroll.

### Real-time with offline capability

The system updates in real-time when connected to the internet. States persist automatically. No user action is needed. Writing and limited access should be possible while offline. Data is synced when online.

### Multi-mode communication

Allow text, audio and video communications. Fall back to asynchronous communication modes when real-time communication fails.

### Interoperability

A comprehensive HTTP API. Allow data import and export. Integrate with existing tools and services wherever possible.

### Ease of use

Provide clean UX with clean, distraction-free, accessibility-friendly UIs. Distraction-free mode is the default. Only show elements relevant to context. Easily switch between views and threads.

### Zero-config setup

 The system should be usable right away after installation.

### Decentralised

No single entity should control the data. Data storage should be decentralised and distributed wherever possible. Built on [Matrix](https://matrix.org).

## Hypothesis

Individuals interact with other individuals over time and in different contexts, and generate information through discussion. In the process, they make use of multiple props and actions to convey their intent and thoughts.

Like the physical world, an individual can watch, participate in, or start a discussion. Depending on their location and the time, the individual will be able to learn a scoped context of the discussion. This is equivalent to asking another individual nearby, "what's going on here?"

The more they are curious the more they get to know how the discussion has reached this far. They can ask other individuals, "Give me the gist".

Unlike the physical world, the individual can then choose to follow their curiosity, and visit all connected discussions in all of recorded space-time. They may not want all the details all the time, and so they can ask for a summary, with a higher level view of the key events in the discussions. They may then choose to focus on specific parts of the discussion, and land in the scoped context of that event in the discussion.

In the flow of a discussion, an individual, or a group of individuals, can have related thoughts that they consider meaningful, but which may end up sidetracking the original discussion. Such related discussions often end up sprouting new thoughts and ideas, and deserve a conversation of their own.

An individual may also notice that a similar discussion has happened earlier, and they may wish to reference the other discussion in the present one. In some cases, they may see it fit to carry on the discussion in the previous group. These lead to multiple pathways that event narratives can take.

In reality, event narratives form complexly interwoven structures. An individual can ask to know the different pathways that have lead to their current vantage point, the different pathways that lead away from there, and the related discussions that sparked from there.

## Data model

### Event

Events are nodes that form the unit of discussion in Upspeak.

- Who said it?
- When did they say it?
- Has this already occurred, or will this occur at a later time?
- Where did they say it?
- What did they say?
- How can I interpret and understand what they said?
- Were they responding to something? If yes, what?
- Where did this discussion begin?
- Can others respond to this? If not, what should they take away from this?

### Context

Context graph ties together Nodes in a meaningful way. Contexts are the edges between event nodes. Think of them as natural conversation threads. The context changes based on the current point of view of the reader.

- How fresh is this discussion?
- What does the whole conversation look like at a glance?
- What conversations lead to here?
- What conversations unfolded from here?
- What related conversations were started from here?
- What did the speaker say before this?
- What did the speaker say after this?
- Who else has mentioned or cited this?
- Who else participated in this discussion before?
- Who else participated in this discussion after?
- What are others talking about based on this?
- Where is fresh discussion happening based on the directions the conversation took from here?
- What other parallel discussions happened near this?
- Which other discussions are relevant to this?
- How does this whole conversation compare to another conversation?
- What should the reader be aware of before reading this conversation?

## License

[Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0)
