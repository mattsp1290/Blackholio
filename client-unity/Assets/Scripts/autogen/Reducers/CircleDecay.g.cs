// THIS FILE IS AUTOMATICALLY GENERATED BY SPACETIMEDB. EDITS TO THIS FILE
// WILL NOT BE SAVED. MODIFY TABLES IN YOUR MODULE SOURCE CODE INSTEAD.

#nullable enable

using System;
using SpacetimeDB.ClientApi;
using System.Collections.Generic;
using System.Runtime.Serialization;

namespace SpacetimeDB.Types
{
    public sealed partial class RemoteReducers : RemoteBase
    {
        public delegate void CircleDecayHandler(ReducerEventContext ctx, CircleDecayTimer timer);
        public event CircleDecayHandler? OnCircleDecay;

        public void CircleDecay(CircleDecayTimer timer)
        {
            conn.InternalCallReducer(new Reducer.CircleDecay(timer), this.SetCallReducerFlags.CircleDecayFlags);
        }

        public bool InvokeCircleDecay(ReducerEventContext ctx, Reducer.CircleDecay args)
        {
            if (OnCircleDecay == null) return false;
            OnCircleDecay(
                ctx,
                args.Timer
            );
            return true;
        }
    }

    public abstract partial class Reducer
    {
        [SpacetimeDB.Type]
        [DataContract]
        public sealed partial class CircleDecay : Reducer, IReducerArgs
        {
            [DataMember(Name = "timer")]
            public CircleDecayTimer Timer;

            public CircleDecay(CircleDecayTimer Timer)
            {
                this.Timer = Timer;
            }

            public CircleDecay()
            {
                this.Timer = new();
            }

            string IReducerArgs.ReducerName => "CircleDecay";
        }
    }

    public sealed partial class SetReducerFlags
    {
        internal CallReducerFlags CircleDecayFlags;
        public void CircleDecay(CallReducerFlags flags) => CircleDecayFlags = flags;
    }
}
