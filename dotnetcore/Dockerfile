FROM mcr.microsoft.com/dotnet/core/aspnet:2.2 AS base
WORKDIR /app
EXPOSE 5001

FROM mcr.microsoft.com/dotnet/core/sdk:2.2 AS build
WORKDIR /src
COPY ["src/postprocessing/postprocessing.csproj", "src/postprocessing/"]
RUN dotnet restore "src/postprocessing/postprocessing.csproj"
COPY . .
WORKDIR "/src/src/postprocessing"
RUN dotnet build "postprocessing.csproj" -c Release -o /app

FROM build AS publish
RUN dotnet publish "postprocessing.csproj" -c Release -o /app

FROM base AS final
WORKDIR /app
COPY --from=publish /app .
LABEL key="jobscheduling-api"
ENTRYPOINT ["dotnet", "postprocessing.dll"]
